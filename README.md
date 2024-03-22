# Scalable-chat-with-redis
A simple but scalable ChatRoom app where users can create/join multiple chat rooms and communicate within them.
Only users in same chat group should receive messages for that group.

How it works:

Redis pubsub channel is created for every group and interested users in group subscribe to that channel.

To connect: ``` ws://localhost:8080/chat?userId="user1" ```

Below commands sent over socket after successful connection

To join channel:
```
{
"cmd":"I-JC",
"chName":"ch1"
}
```

To send message to channel:
```
{
"cmd":"I-SM",
"chName":"ch1",
"msg":"i dont know"
}
```
