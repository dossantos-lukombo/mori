document.addEventListener("DOMContentLoaded", function() {
    const chatWindow = document.getElementById("chat-window");
    const userInput = document.getElementById("user-input");
    const sendButton = document.getElementById("send-button");
  
    function appendMessage(sender, message) {
      const messageElement = document.createElement("div");
      messageElement.classList.add("chat-message", sender);
  
      const messageContent = document.createElement("div");
      messageContent.classList.add("message-content");
      messageContent.textContent = message;
  
      messageElement.appendChild(messageContent);
      chatWindow.appendChild(messageElement);
      chatWindow.scrollTop = chatWindow.scrollHeight;
    }
  
    function handleUserMessage() {
      const message = userInput.value.trim();
      if (message === "") return;
  
      appendMessage("user", message);
      userInput.value = "";
  
      // Simuler une réponse du chatbot
      setTimeout(() => {
        appendMessage("bot", "Ceci est une réponse simulée du chatbot.");
      }, 1000);
    }
  
    sendButton.addEventListener("click", handleUserMessage);
  
    userInput.addEventListener("keypress", function(event) {
      if (event.key === "Enter") {
        handleUserMessage();
      }
    });
  });
  