> 尽量使用root账户
# JAVA
```shell
# SH
sudo ./autostart add myapp "sh /home/software/jar/omniscient_test/run2.sh start -b false" --workdir=/home/software/jar/omniscient_test

# JDK 
## 简单
sudo ./autostart add myapp "java -jar /home/software/jar/omniscient_test/omniscient_test-0.0.1-SNAPSHOT.jar" --workdir=/home/software/jar/omniscient_test
## 参数丰富
sudo ./autostart add myapp "java -Xms1024m -Xmx4024m -Xmn4024m -Dfile.encoding=UTF-8 -jar /home/software/jar/omniscient_test/omniscient_test-0.0.1-SNAPSHOT.jar --server.port=5122 --spring.profiles.active=prod" --workdir=/home/software/jar/omniscient_test
```


# Python
> 我测试的py3的脚本
```python
import time

# 获取当前时间并格式化为指定格式
current_time = time.strftime("%Y%m%d%H%M%S", time.localtime())

# 使用格式化后的时间作为文件名
file_name = f"{current_time}-log.txt"

with open(file_name, "w", encoding="utf-8") as wp:
    line = f"{current_time} 開始=========="
    wp.write(line)
```
```shell
autostart add testPy "python3 /home/software/py/test_hello.py"
```