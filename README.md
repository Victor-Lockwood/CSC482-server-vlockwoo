# Server

## How To Deploy

1. Kick off the multistage build with `docker build -t vlockwoo-server .`
2. To test locally, run `docker run -p 40000:8090 -d server` then in a 
browser, hit `localhost:40000/vlockwoo/status`.  You should get JSON with a 
200 response back.  Feel free to kill the local container.
3. In the same directory as your pem for authentication 
(mainly for convenience’s sake), run `docker save --output vlockwoo-server.tar server`
4. Run `scp -i <pem file name> vlockwoo-server.tar <ec2 username>@<ec2 IP>:`
5. SSH into the EC2 instance
6. Run `docker load --input vlockwoo-server.tar`
7. Run `docker run -e LOGGLY_TOKEN=<token> -p 40000:8090 -d vlockwoo-server`
8. Hit `http://<ip>:40000/vlockwoo/status` and if all is well you should have gotten a 200 response with JSON.
9. Clean up and remove your `.tar` file on the EC2 instance with `rm vlockwoo-server.tar`.
10. You can now exit the SSH session.