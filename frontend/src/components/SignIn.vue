<template>

  <div class="sign-in__wrapper">
    <div class="image-div-login">
      <img src="../assets/toa-heftiba-l_ExpFwwOEg-unsplash.jpg" alt="people hanging out">

    </div>

    <div class="sign-in">
      <h1 class="mori">Mori <span class="adder">- sign in</span></h1>
      <form class="form-group" @submit.prevent="signSubmit" id="sign-in__form">
        <div class="form-input">
          <label for="username">Email</label>
          <input type="email" id="email" v-model="signInForm.login" required>
        </div>
        <div class="form-input">
          <label for="password">Password</label>
          <input type="password" id="password" v-model="signInForm.password" required>
        </div>
      </form>
      <div class="button-or-register">
        <button class="btn" form="sign-in__form" type="submit">Sign in</button>
        <p>Need an account?
          <router-link to="/reg" id="sign-up">Register here</router-link>
        </p>
      </div>
    </div>
  </div>

</template>


<script>
export default {
  name: "SignIn",
  data() {
    return {
      signInForm: {
        login: "",
        password: "",
      },
    };
  },
  methods: {
    toast() {
      /*---------------           Here is toast example             --------------------*/
      // 
      this.$toast.open({
        message: "Data sent!",
        type: "default",
        //optional options
        position: "bottom-right",
        duration: 3000,
        dismissible: true,
        onClick: null,
        onDismiss: null,
        queue: false,
        pauseOnHover: true, //Pause the timer when mouse on over a toast
      });
    },
    async signSubmit() {
      try {
        // await fetch('https://bfdf8b79-b1e1-40ce-8d02-896de58da3ca.mock.pstmn.io/signin', {
        await fetch("http://localhost:8081/signin", {
          credentials: "include",
          method: "POST",
          headers: {
            "Accept": "application/json",
            "Content-Type": "application/json"
          },
          body: JSON.stringify(this.signInForm)
        })
          .then((response => response.json()))
          .then((json => {
            // console.log(json)
            if (json.message === "Login successful") {
              this.$toast.open({
                message: "Login success!",
                type: "success", //One of success, info, warning, error, default
              });


              this.$store.dispatch("createWebSocketConn").then(() => this.$router.push("/main"))


              // await this.$store.dispatch("getMyUserID")
            }
            else {
              this.$router.push("/");
              this.$toast.open({
                message: json.message,
                type: "error", //One of success, info, warning, error, default
              });
            }
          }));
      }
      catch { }
    },
  }
}
</script>

<style >



.sign-in__wrapper {
  display: flex;
  /* margin: auto 0; */
  background-color: var(--bg-neutral);
  border-radius: 20px;
  color: var(--color-white);
  box-shadow: 0 4px 15px rgb(0, 0, 0);
  overflow: hidden;
  align-items: center;
  width: 810px;
}

.image-div-login {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50%;
  height: 100%;
  background-color: var(--bg-neutral);
}

.sign-in__wrapper img {
  height: 550px;
  background-position: center;
  background-repeat: no-repeat;
  width: 100%;
}


.sign-in {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 40px;
  margin: 0 auto;
  padding: 0 50px;
  width: 60%;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 15px;
  color: var(--color-white);
  width: 90%;
}

.form-input {
  color: var(--color-white);
}



.sign-in button {
  margin-bottom: 10px;
  width: 80px;
}



#sign-up {
  color: var(--purple-color);
  text-decoration: underline;
  font-size: 16.5px;  
  transition: all 0.3s ease;
}

#sign-up:hover {
  color: var(--hover-background-color);
  text-decoration: underline;
  font-size: 16.5px;
}

.button-or-register {
  display: flex;
  flex-direction: column;
  gap: 20px;
}


</style>