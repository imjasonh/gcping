# Install Python
sudo apt-get install software-properties-common python-software-properties
sudo add-apt-repository ppa:fkrull/deadsnakes-python2.7
sudo apt-get update
sudo apt-get install python2.7

# Start simple HTTP server with CORS header set.
echo "pong" > ping
cat << EOF > cors-http-server.py
#! /usr/bin/env python

from SimpleHTTPServer import SimpleHTTPRequestHandler, test

class CORSHTTPRequestHandler(SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        super(CORSHTTPRequestHandler, self).end_headers(self)

if __name__ == '__main__':
    test(HandlerClass=CORSHTTPRequestHandler)
EOF
sudo python cors-http-server.py