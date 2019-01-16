import Vue from 'vue'
import VueRouter from 'vue-router'
import $auth from "./components/auth";
import App from './components/App';
import Login from './components/Login';
import NewAccount from './components/NewAccount';
import Dashboard from './components/Dashboard';
import * as VueGoogleMaps from 'vue2-google-maps';
import VModal from 'vue-js-modal';

Vue.use(VueRouter);
Vue.use(VueGoogleMaps, {
    load: {
        key: 'AIzaSyDf43lPdwlF98RCBsJOFNKOkoEjkwxb5Sc'
    }
});
Vue.use(VModal);

function requireAuth(to, from, next) {
    if (!$auth.isLoggedIn()) {
        next({
            path: '/login',
            query: {redirect: to.fullPath}
        })
    } else {
        next()
    }
}

const router = new VueRouter({
    mode: 'history',
    routes: [
        {path: '/login', component: Login},
        {path: '/new', component: NewAccount},
        {path: '/dashboard', component: Dashboard, alias: '/', beforeEnter: requireAuth},
        {
            path: '/logout', beforeEnter(to, from, next) {
                $auth.logout();
                next('/login');
            }
        },
    ]
});

new Vue({
    el: '#app',
    router,
    render: h => h(App),
});
