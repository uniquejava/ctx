kubectl切换context和默认的namespace好复杂,能不能写一个golang程序, 包含下面的命令(子命令)
1. ctx ls列出所有的context(应该是读取~/.kube/config)
2. 当前的context前面标上*号
3. ctx use xxx就是切换context
4. ctx use xxx default-ns, 可以同时设置默认的ns
5. ctx rm xxx就是删除context
6. 如果可以, 能用上下箭头和Enter键选择context更好