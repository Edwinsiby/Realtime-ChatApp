// script.js

const chatMessages = document.getElementById('chat-messages');
const messageInput = document.getElementById('message-input');
const sendButton = document.getElementById('send-button');
const loginForm = document.querySelector('form');

const ws = new WebSocket('ws://localhost:8080/ws/chat');
let userName = ''; 

ws.onopen = () => {
    console.log('WebSocket connected');
};

ws.onmessage = (event) => {
    const receivedMessage = event.data;
    appendMessage(receivedMessage);
};

ws.onclose = () => {
    console.log('WebSocket disconnected');
};

sendButton.addEventListener('click', () => {
    const message = messageInput.value;
    if (message) {
        ws.send(message);
        messageInput.value = '';
    }
});

loginForm.addEventListener('submit', (e) => {
    e.preventDefault();
    userName = getUsernameFromInput(); 
    if (userName) {
        ws.send(userName);
        updateWelcomeMessage(userName);
        clearUsernameInput();
    }
});

function getUsernameFromInput() {
    const usernameInput = document.querySelector('input[name="userName"]');
    return usernameInput.value;
}

function updateWelcomeMessage(username) {
    const welcomeMessage = document.getElementById('user-name');
    welcomeMessage.textContent = username;
}

function clearUsernameInput() {
    const usernameInput = document.querySelector('input[name="userName"]');
    usernameInput.value = '';
}

function appendMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.textContent = message;
    chatMessages.appendChild(messageElement);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}
