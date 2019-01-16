<template>
    <div class="wrapper">
        <h2>Login</h2>
        <form @submit.prevent="login" class="login">
            <input placeholder="username" v-model="username" type="text">
            <input placeholder="password" v-model="password" type="password">
            <input type="submit" value="Login">
        </form>
        <div>Don have an account?
            <router-link to="/new">Create one</router-link>
        </div>

        <div class="error">
            {{ errorMsg }}
        </div>
    </div>
</template>

<script>
    import $auth from './auth.js'

    export default {
        name: 'Login',
        data() {
            return {
                username: '',
                password: '',
                errorMsg: null
            }
        },

        methods: {
            login() {
                $auth.login(this.username, this.password).then(response => {
                    localStorage.token = response.data.token;
                    this.$router.replace(this.$route.query.redirect || '/dashboard');
                }).catch(() => {
                    this.errorMsg = "Wrong username/password";
                });
            }
        }

    }
</script>

<style scoped>
    .wrapper {
        display: grid;
        justify-content: center;
        margin-top: 3rem;
        text-align: center;
    }

    .login {
        row-gap: .3rem;
        display: grid;
    }

    .login > * {
        font-size: 1rem;
        height: 2rem;
    }

    .error {
        color: red;
    }
</style>
