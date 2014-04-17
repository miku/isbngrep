isbngrep: isbngrep.go
	gofmt -w -tabs=false -tabwidth=4 isbngrep.go
	go build -o isbngrep isbngrep.go

clean:
	rm -f isbngrep
	rm -f isbngrepc

# buildrpm: https://gist.github.com/miku/7874111
rpm: isbngrep
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp isbngrep.spec $(HOME)/rpmbuild/SPECS
	cp isbngrep $(HOME)/rpmbuild/BUILD
	./buildrpm.sh isbngrep
	cp $(HOME)/rpmbuild/RPMS/x86_64/*rpm .

vagrant.key:
	curl -sL "https://raw2.github.com/mitchellh/vagrant/master/keys/vagrant" > vagrant.key
	chmod 0600 vagrant.key

# helper to build RPM on a RHEL6 VM, to link against glibc 2.12
# Assumes a RHEL6 go installation (http://nareshv.blogspot.de/2013/08/installing-go-lang-11-on-centos-64-64.html)
# And: sudo yum install git rpm-build
# Don't forget to vagrant up :)
rpm-compatible: vagrant.key
	ssh -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1 -p 2222 "mkdir -p /home/vagrant/github/miku && cd /home/vagrant/github/miku && [ -d /home/vagrant/github/miku/isbngrep ] || git clone https://github.com/miku/isbngrep.git"
	ssh -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1 -p 2222 "cd /home/vagrant/github/miku/isbngrep && git pull origin master && GOPATH=/home/vagrant go get github.com/codegangsta/cli && GOPATH=/home/vagrant make rpm"
	scp -o port=2222 -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1:/home/vagrant/github/miku/isbngrep/*rpm .

vm-setup:
	ssh -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1 -p 2222 "sudo yum install -y http://ftp.riken.jp/Linux/fedora/epel/6/i386/epel-release-6-8.noarch.rpm"
	ssh -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1 -p 2222 "cd /etc/yum.repos.d && sudo wget http://www.hop5.in/yum/el6/hop5.repo"
	ssh -o StrictHostKeyChecking=no -i vagrant.key vagrant@127.0.0.1 -p 2222 "sudo yum install -y golang git rpm-build"
