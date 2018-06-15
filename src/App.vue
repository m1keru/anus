<template>
  <div id="app">
    <AnsibleScripts/>
  </div>
</template>

<script>
import AnsibleScripts from './components/AnsibleScripts'

export default {
  name: 'App',
  components: {
    AnsibleScripts
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 0px;
}
</style>
