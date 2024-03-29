# snowFlake 雪花算法
snowflake ID 算法是 twitter 使用的唯一 ID 生成算法，为了满足 Twitter 每秒上万条消息的请求，使每条消息有唯一、有一定顺序的 ID ，且支持分布式生成。

## Install 

`go get -u /github.com/trumpliu2018/snowflake`

## 原理
某一台拥有独立标识(为机器分配独立id)的机器在1毫秒内生成带有不同序号的id 
所以生成出来的id是具有时序性和唯一性的

- 构成

snowflake ID 的结构是一个 64 bit 的 int 型数据。

41 bit 时间戳 + 10bit工作机器 + 12bit 序列号.

- **第1位bit：** 

二进制中最高位为1的都是负数，但是我们所需要的id应该都是整数，所以这里最高位应该为0

- **后面的41位bit：** 
用来记录生成id时的毫秒时间戳，这里毫秒只用来表示正整数(计算机中正整数包含0)，所以可以表示的数值范围是0至2^41 - 1（这里为什么要-1很多人会范迷糊，要记住，计算机中数值都是从0开始计算而不是1， 41bit大约可以用68年

- **再后面的10位bit：**
用来记录工作机器的id 
2^10 = 1024 所以当前规则允许分布式最大节点数为1024个节点 我们可以根据业务需求来具体分配worker数和每台机器1毫秒可生成的id序号number数

- **最后的12位：** 
用来表示单台机器每毫秒生成的id序号 
12位bit可以表示的最大正整数为2^12 - 1 = 4096，即可用0、1、2、3...4095这4096(注意是从0开始计算)个数字来表示1毫秒内机器生成的序号(这个算法限定单台机器1毫秒内最多生成4096个id，超出则等待下一毫秒再生成)
最后将上述4段bit通过位运算拼接起来组成64位bit

## 改进
由于前段js 中 number最长位53位, 为了保证js数值不被截断, 机器位使用8bit, 序列号使用6bit
这样可以保证256个机器. 每毫秒生产64个序列号.

## 使用
```
func TestSnowFlake(t *testing.T) {
	worker, err := NewNode(1)

	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan int64)
	count := 100000
	for i := 0; i < count; i++ {
		go func() {
			id := worker.Generate()
			println(id)
			ch <- id
		}()
	}

	// defer close(ch)

	m := make(map[int64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		_, ok := m[id]
		if ok {
			t.Error("ID is not unique!")
			return
		}
		m[id] = i
	}
	fmt.Println("All", count, " successed!")
}
```
