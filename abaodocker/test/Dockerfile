FROM    ccr.ccs.tencentyun.com/abao/centos

RUN     /bin/echo 'root:123456' |chpasswd
RUN     useradd runoob
RUN     /bin/echo 'runoob:123456' |chpasswd
RUN     /bin/echo -e "LANG=\"en_US.UTF-8\"" >/etc/default/local && \
pwd
EXPOSE  22
EXPOSE  80
CMD     /usr/sbin/sshd -D