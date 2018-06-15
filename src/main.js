// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import 'bootstrap'
import 'jquery'
import Vue from 'vue'
import App from './App'
import VueResource from 'vue-resource'
import './../node_modules/jquery/dist/jquery.min.js';
import './../node_modules/bootstrap/dist/css/bootstrap.min.css';
import './../node_modules/bootstrap/dist/js/bootstrap.min.js';
import VueFontAwesomeCss from 'vue-fontawesome-css'
Vue.use(VueResource)
Vue.use(VueFontAwesomeCss)

Vue.config.productionTip = false
// Vue.config.devtools = true

/* eslint-disable no-new */
new Vue({
  el: '#app',
  components: { App },
  template: '<App/>',
  data: {
    loading: false
  }
})
