// import { marked } from "https://cdn.jsdelivr.net/npm/marked/lib/marked.esm.js";
// import markdownToTxt from 'https://cdn.jsdelivr.net/npm/markdown-to-txt@2.0.1/+esm'// import markdownToTxt from 'markdown-to-txt';

const chatBox = document.querySelector(".chat-window");

let conversation = {
    conversation_id: "",
    user_request: "",
    llm_response: "",
    session:"",
    new_conversation: false,
    history: []
}
// Fonction pour afficher un message dans la chat-box
function appendMessage(sender, message) {
    const messageElement = document.createElement("pre");
    messageElement.classList.add("message", sender);

    messageElement.innerText = message
    
    chatBox.appendChild(messageElement);
    chatBox.scrollTop = chatBox.scrollHeight;  // Scroll automatiquement vers le bas
}

// Fonction pour envoyer des données au serveur Go
function sendBtn_clicked(){
    // let userID = window.location.href.split("/").pop()
    let userID = "0001"
    let buttonSend = document.getElementById("send-btn")
    let textearea = document.getElementById("user-input")
    buttonSend.addEventListener("click", (e)=>{
        e.preventDefault()
        console.log("btn clicked !")

        conversation.user_request = textearea.value
        appendMessage("Utilisateur", textearea.value)
        sendData(userID,textearea.value)
        textearea.value = ""
        
    })
    textearea.addEventListener("keydown", (e)=>{
        
        if(e.key === "Enter" && e.shiftKey){
            e.preventDefault()
            console.log("Shift + Enter clicked !")
            textearea.value += "\n"
        }else if(e.key === "Enter"){
            e.preventDefault()
            console.log("Enter clicked !")

            conversation.user_request = textearea.value
            appendMessage("Utilisateur", textearea.value)
            sendData(userID,textearea.value)
            textearea.value = ""
        }
    })

    return null
}

async function sendData(userID,data) {

    console.log("session: ",userID)
    console.log("data: ", data)

    let result = null
    const messageElement = document.createElement("pre");
    
    const response = await fetch("http://localhost:8080/conversation/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",            
        },
        body: JSON.stringify({ session: userID, message: data })
    });

    if (!response.ok || !response.body) {
        console.error("Erreur lors de l'envoi des données :", response.statusText);
        return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");

    try {
        let { done, value } = await reader.read();

        while (!done) {
            // Convertir le chunk en texte
            const chunk = decoder.decode(value, { stream: true });
            
            messageElement.classList.add("message", "LLM");

            // Diviser les événements en fonction du format SSE
            const lines = chunk.split("\n");
            for (const line of lines) {
                if (line.startsWith("data: ")) {
                    console.log("line: ",line)
                    let jsonData = line.replace("data: ", "");
                    // let jsonData = line.split(" ")
                    jsonData = jsonData.trim()
                    try {
                        console.log(jsonData)
                        let parsedData = JSON.parse(jsonData);
                        // responseMessage = parsedData.response; 
                        messageElement.innerText += parsedData.response;
                    } catch (e) {
                        console.error("Erreur de parsing JSON :", e);
                    }
                }
            }
            // appendMessage("LLM", responseMessage); // Afficher chaque chunk
            
            // Lire le prochain chunk
            ({ done, value } = await reader.read());
            chatBox.appendChild(messageElement);
            chatBox.scrollTop = chatBox.scrollHeight;
        }
        conversation.llm_response = messageElement.innerText
        sendConversation()
    } catch (error) {
        console.error("Erreur de lecture du flux", error);
    } finally {
        reader.releaseLock();
    }

}

async function sendConversation(){

    conversation.session = window.location.href.split("/").pop()

    console.log("CONVERSATION: ",conversation)

    const response = await fetch("http://localhost:8080/api/conversation", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(conversation)
    });

    if (!response.ok || !response.body) {
        console.error("Erreur lors de l'envoi des données :", response.statusText);
        return;
    }else{
        console.log("Conversation envoyée avec succès !")

        // const reader = response.json();
        // console.log("reader server response : ",reader)
    }
}

function SendNewChat(){
    const session = window.location.href.split("/").pop()
    const formNewChat = document.getElementById("form-new-chat")
    const btnNewChat = document.getElementById("btn-new-chat")
    formNewChat.setAttribute("action",`/home/${session}`)
    btnNewChat.addEventListener("click", (e)=>{
        e.preventDefault()
        
    })
    formNewChat.addEventListener("submit", (e)=>{
        e.preventDefault()
    })
    buttonSend.addEventListener("click", (e)=>{
        e.preventDefault()
        console.log("btn clicked !")

        conversation.user_request = textearea.value
        appendMessage("Utilisateur", textearea.value)
        sendData(userID,textearea.value)
        textearea.value = ""
    })

    return null
}



sendBtn_clicked()


export {sendBtn_clicked}