INSERT INTO `jpid` (`name`, `port`, `pid`, `catalog`, `run`, `description`) VALUES
                                                                                ('UserManage', 8080, 123, '/usr/local/java/usermanage', 'sh /usr/local/java/usermanage/start.sh', '用户管理系统'),
                                                                                ('OrderSys', 8081, 124, '/opt/java/ordersys', 'sh /opt/java/ordersys/run.sh', '订单处理系统'),
                                                                                ('ProductSys', 8082, 125, '/home/java/productsys', 'sh /home/java/productsys/init.sh', '商品信息管理系统'),
                                                                                ('AuthSys', 8083, 126, '/var/lib/java/authsys', 'sh /var/lib/java/authsys/launch.sh', '认证授权系统'),
                                                                                ('ReportSys', 8084, 127, '/root/java/reportsys', 'sh /root/java/reportsys/generate_report.sh', '报表生成系统');
