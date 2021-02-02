var mqtt=require('mqtt');
const fs = require('fs');
var ca = fs.readFileSync('/home/dusan/go/src/github.com/mainflux/mainflux/docker/ssl/certs/ca.crt');
var key = fs.readFileSync('/home/dusan/go/src/github.com/mainflux/mainflux/docker/ssl/certs/thing.key');
var cert = fs.readFileSync('/home/dusan/go/src/github.com/mainflux/mainflux/docker/ssl/certs/thing.crt');
var connOptions = {
    clientId:"ivkenodejs",
    ca: ca,
    key: key,
    cert: cert,
    // rejectUnauthorized : false,
    protocol: 'wss',
    username:"steve",
    password:"password",
};
var msgoptions={
    retain:true,
    qos:2};
var client  = mqtt.connect("wss://localhost:8081/mqtt", connOptions);
console.log("connected flag:  " + client.connected);
client.on("connect",function(){	
    console.log("connected flag:  "+client.connected);
    console.log("connected");
    client.publish("testtopic", "test message", msgoptions)
    // setInterval(function(){
    //     client.publish("testtopic", "test message", msgoptions);
    // }, 3000);
    // client.end();
    // console.log("connected flag:  "+client.connected);
});
client.on("error",function(error){
    console.log("Can't connect" + error);
    process.exit(1)}
);
client.subscribe("testtopic",{qos:1});
client.on('message',function(topic, message, packet){
	console.log("message is: "+ message);
	console.log("topic is: "+ topic);
});