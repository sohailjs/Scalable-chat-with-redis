// App.js
import React, {useState} from 'react';
import './ChatApp.css'

const connectURL=`ws://localhost:8080/chat?userId=`
const App = () => {
    const [username, setUsername] = useState('');
    const [enteredChatRoom, setEnteredChatRoom] = useState(false);

    const handleUsernameChange = (event) => {
        setUsername(event.target.value);
    };

    const handleEnterChatRoomClick = () => {
        if (username.trim() !== '') {
            const ws = new WebSocket(`${connectURL}${username}`);
            ws.onopen = () => {
                // Send the username to the server upon connection
                ws.send(JSON.stringify({type: 'username', username: username}));
                setEnteredChatRoom(true);
            };

            // Handle messages received from the WebSocket server
            ws.onmessage = (message) => {
                console.log('Received message:', message.data);
                // Handle messages as needed
            };

            // Handle WebSocket errors
            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
        } else {
            alert('Please enter a username first!');
        }
    };

    return (
        <div>
            <input
                type="text"
                placeholder="Enter your username"
                value={username}
                onChange={handleUsernameChange}
            />
            <button onClick={handleEnterChatRoomClick}>Enter Chat Room</button>
            {enteredChatRoom && <p>Entered chat room as {username}</p>}
        </div>
    );
};

export default App;
