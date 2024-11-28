import ollama


historique = [{"role": "system", "content": "Vous êtes un assistant utile."}]

def treating_user_request(output:dict):
    client = ollama.Client(host="http://localhost:11434")
    
    historique.append({"role": "user", "content": output["message"]})
    
    response = client.chat(
    model="llama3.2:3b",
    messages=historique,
    stream=True,
    )
    stream = ""
    print("response: ", response)
    for chunk in response:
        print("\nchunk: \n",chunk)
        
        c = chunk["message"]["content"]
        # print(c, end='', flush=True)
        stream += c
        
        yield chunk
    # return response
    
    historique.append({"role": "assistant", "content": stream})

def conversation_resumed(output:dict):
    
    client = ollama.Client(host="http://localhost:11434")

    
    historique.append({"role": "assistant", "content": output["message"]})
    historique.append({"role": "user", "content": "résume moi ce que nous avons dit précédemment une phrase de maximum 5 mots dans la langue de la conversation."}) 
    
    response = client.chat(
    model="llama3.2:3b",
    messages=historique,
    )
    
    response = response.choices[0].message.content

    
    
    