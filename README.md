  
![](./img/logo.png)
- 思想的萌芽🌱：  
编写这个工具最初的目的是因为，使用企业账号收集资产，权益积分很多，但是导出额度没有，导致我在有很多查询积分的情况下，没办法调用Hunter给的全部查询并导出结果为文件的api接口，所以就只能分页遍历查询并导出为excel文件，反正写都写了，就索性把该有的查询一起写了，当作以后工作中的工具来用。  
- 不断的进步💪：  
如果真的有其他人也会用我的工具，就真的太开心了，师傅们如果发现Bug或者感觉哪里需要改进，告诉我，我会更加的努力。
- 暂时的不足💻：  
1.因为工作时间原因，第一版生成运行文件体积较大，回头抽时间优化一下 2.最初用多线程跑查询，发现hunter的接口好像调用有限制，甚至单线程跑的快了都不行，所以循环遍历调用接口的时候，做了3秒到等待，导致速度可能比较慢，后面再详细查一下具体原因是因为hunter的接口调用做了频率限制还是其他原因。
## 使用方法
需要在和软件的同一目录位置创建配置文件：hunterx.yaml
文件内容：(yaml文件格式，记得冒号后面打一个空格哦)
```yaml
#hunter_userName
userName: 
#hunter_apiKey
apiKey: 
```
- 参数：
![](./img/-h.png)
- 事例：  
`HunterX -q 'domain="g.mi.com"'`  
![](./img/-q.png)  
`HunterX -q 'domain.suffix="xx.com" or domain.suffix="xx.com"' -all true`  
![](./img/-all.png)
`HunterX -l target.txt`  
![](./img/-l.png)