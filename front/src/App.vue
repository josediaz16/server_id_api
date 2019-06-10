<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <HelloWorld msg="Welcome to Server ID Service"/>
    <SearchBox plh="Type a server domain"/>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import SearchBox  from './components/SearchBox.vue'

export default {
  name: 'app',
  components: {
    HelloWorld,
    SearchBox
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
  color: #2c3e50;
  margin-top: 60px;

}

#app img {
  width: 300px;
}
</style>
