document.addEventListener("DOMContentLoaded", function () {
  const chat_window = document.querySelector(".chat-window");
  const welcome_message = document.createElement("div");
  welcome_message.innerText = "Bienvenue chez Mori";
  welcome_message.classList.add("welcome-message");
  welcome_message.classList.add("full-width");
  welcome_message.style.textAlign = "center";
  welcome_message.style.fontWeight = "bold";
  welcome_message.style.fontSize = "20px";
  welcome_message.style.display = "block";
  chat_window.appendChild(welcome_message);

  const buttonSend = document.getElementById("send-btn");

  buttonSend.addEventListener("click", function () {
    welcome_message.scrollIntoView();
    welcome_message.remove();
  });
});

