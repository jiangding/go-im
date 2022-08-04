package controller

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)
const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)
type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"` //消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"` //谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`//对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"` //消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"` //预览图片
	Url     string `json:"url,omitempty" form:"url"` //服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"` //简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"` //其他和数字相关的
}
/*
消息发送结构体
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}
2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,content:"标题",pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/dsturl","memo":"这是描述"}
3、MEDIA_TYPE_VOICE，amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,url:"http://www.baidu.com/a/log,jpg"}
5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}
7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/a.mp4"}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,"content":"10086","pic":"http://www.baidu.com/a/avatar,jpg","memo":"胡大力"}
*/


// Node 本核心在于形成userid和Node的映射关系
type Node struct {
	Conn *websocket.Conn
	//并行转串行
	DataQueue chan []byte // 通道
	GroupSets set.Interface // 集合
}

// 映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node,0)
// 读写锁
var rwlocker sync.RWMutex

// Chat ws://127.0.0.1/chat?id=1&token=xxxx
func Chat(w http.ResponseWriter, req *http.Request){

	// 检验接入是否合法
	query := req.URL.Query()
	id := query.Get("id")
	token := query.Get("token")

	// 字符串转整型
	userId,_ := strconv.ParseInt(id, 10, 64)

	// 检查token
	isValida := checkToken(userId, token)

	conn,err :=(&websocket.Upgrader{
		// 客户端发送 HTTP 请求，请求服务器将用于 HTTP 请求的连接升级为 WebSocket 协议
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(w,req,nil)
	if err!=nil{
		log.Println(err.Error())
		return
	}

	// 获得conn
	node := &Node{
		Conn: conn,
		DataQueue: make(chan []byte, 50),
		GroupSets:set.New(set.ThreadSafe),
	}

	//todo 获取用户全部群Id
	comIds := contactService.SearchComunityIds(userId)
	for _,v:=range comIds {
		node.GroupSets.Add(v) // 加入 groupSets集合中
	}


	rwlocker.Lock() // 加锁, 保持唯一性处理,不被并行影响
	clientMap[userId] = node   // userid和node形成绑定关系
	rwlocker.Unlock() // 解锁

	// 发送逻辑
	go sendProc(node)
	// 接受逻辑
	go recvProc(node)

	log.Printf("<-%d\n",userId)
	// 发送消息
	sendMsg(userId, []byte("hello world!"))
}

//todo 发送消息
func sendMsg(userId int64,msg []byte) {
	rwlocker.RLock() // 加读取锁,锁越细分效率越高
	node, ok:= clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg  // 此时写入通道
	}
}

//AddGroupId 添加新的群ID到用户的groupset中
func AddGroupId(userId,gid int64){
	//取得node
	rwlocker.Lock()
	node,ok := clientMap[userId]
	if ok {
		node.GroupSets.Add(gid) // 添加gid到set
	}
	//clientMap[userId] = node
	rwlocker.Unlock()
}

//todo 发送协程
func sendProc(node *Node) {
	for {
		select {
		case data:= <-node.DataQueue: // 监听通道, 获取队列中的data
			err := node.Conn.WriteMessage(websocket.TextMessage,data)  // 发送消息到 WebSocket
			if err!=nil{
				log.Println(err.Error())
				return
			}
		}
	}
}
//todo 接收协程
func recvProc(node *Node) {
	for {
		_,data,err := node.Conn.ReadMessage() // 循环读取 WebSocket 里面的信息
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//dispatch(data)
		//把消息广播到局域网
		broadMsg(data)
		log.Printf("[ws]<=%s\n",data)
	}
}

func init(){
	// 初始化开启
	go udpsendproc()
	go udprecvproc()
}

//用来存放发送的要广播的数据
var  udpsendchan chan []byte = make(chan []byte,1024)
//todo 将消息广播到局域网
func broadMsg(data []byte){
	udpsendchan <- data
}

//todo 完成udp数据的发送协程
func udpsendproc(){
	log.Println("start udpsendproc")
	//todo 使用udp协议拨号
	con,err:=net.DialUDP("udp",nil,
		&net.UDPAddr{
			IP:net.IPv4(192,168,0,255),
			Port:3000,
		})
	defer con.Close()
	if err!=nil{
		log.Println(err.Error())
		return
	}
	//todo 通过的到的con发送消息
	//con.Write()
	for{
		select {
		case data := <- udpsendchan:
			_,err=con.Write(data)
			if err!=nil{
				log.Println(err.Error())
				return
			}
		}
	}
}
//todo 完成upd接收并处理功能
func udprecvproc(){
	log.Println("start udprecvproc")
	//todo 监听udp广播端口
	con,err:=net.ListenUDP("udp",&net.UDPAddr{
		IP:net.IPv4zero,
		Port:3000,
	})
	defer con.Close()
	if err!=nil{log.Println(err.Error())}
	//TODO 处理端口发过来的数据
	for{
		var buf [512]byte
		n,err:=con.Read(buf[0:])
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//直接数据处理
		dispatch(buf[0:n])
	}
	log.Println("stop updrecvproc")
}

// 后端调度逻辑处理,分发数据出去
func dispatch(data []byte){
	//todo 解析data为message
	msg := &Message{}
	err := json.Unmarshal(data,&msg)
	if err != nil{
		log.Println(err.Error())
		return
	}
	//todo 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid,data)
	case CMD_ROOM_MSG:
		//todo 群聊转发逻辑
	case CMD_HEART:
		//todo 一般啥都不做
	}
}

//检测是否有效
func checkToken(userId int64,token string) bool {
	//从数据库里面查询并比对
	user := userService.Find(userId)
	return user.Token==token
}