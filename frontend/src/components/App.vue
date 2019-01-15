<template>
    <div id="app">
        <template>
            <router-view></router-view>
        </template>
    </div>
</template>

<script>
    export default {
        name: 'app',
    }
</script>

<style>
    body {
        color: #404040;
        font: 400 15px/22px 'Source Sans Pro', 'Helvetica Neue', Sans-serif;
        margin: 0;
        padding: 0;
    }

    * {
        box-sizing: border-box;
    }

</style>
