# Trello-Clone
dont forget ssl commands
openssl genpkey -algorithm RSA -out server.key

openssl req -new -key server.key -out server.csr -subj "/C=RS/ST=Serbia/L=Belgrade/O=MyCompany/OU=IT/CN=example.com/emailAddress=admin@example.com"

openssl x509 -req -in server.csr -signkey server.key -out server.crt
