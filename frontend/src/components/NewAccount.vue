<template>
    <div class="wrapper">
        <h2>New Client Account</h2>
        <form @submit.prevent="newAccount" class="newAccount">
            <input placeholder="username" required v-model="username" type="text">
            <input placeholder="password" required v-model="password" type="password">
            <input placeholder="password (confirm)" required v-model="passwordVis" type="password">
            <input type="submit" value="Create Account">
        </form>
        <div class="error">
            {{ errorMsg }}
        </div>
    </div>
</template>


<script>
    import $users from './users';

    export default {
        data() {
            return {
                username: null,
                password: null,
                passwordVis: null,
                errorMsg: null
            }
        },

        methods: {
            newAccount() {
                if (this.password !== this.passwordVis) {
                    this.errorMsg = "Passwords must match";
                    return
                }

                $users.createClientAccount(this.username, this.password).then(res => {
                    alert(`User ${res.data.username} created`);
                    this.$router.replace('/login');
                }).catch(err => {
                    this.errorMsg = err.response.data;
                })
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

    .newAccount {
        row-gap: .3rem;
        display: grid;
    }

    .newAccount > * {
        font-size: 1rem;
        height: 2rem;
    }

    .error {
        max-width: 20rem;
        color: red;
    }
</style>
