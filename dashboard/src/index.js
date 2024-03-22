import React from 'react';
import ReactDOM from 'react-dom';
import ChatApp from './App'; // Assuming your main component is named ChatApp and is located in the same directory

ReactDOM.render(
    <React.StrictMode>
        <ChatApp />
    </React.StrictMode>,
    document.getElementById('root')
);
