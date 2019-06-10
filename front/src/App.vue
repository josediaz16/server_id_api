<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <HelloWorld msg="Welcome to Server ID Service"/>
    <SearchBox
      plh="Type a server domain"
      v-on:search="searchServer"
    />
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import SearchBox  from './components/SearchBox.vue'
import axios      from 'axios'

export default {
  name: 'app',
  components: {
    HelloWorld,
    SearchBox
  },
  methods: {
    searchServer: function(query) {
      axios.get("http://localhost:3000/domains/search", {params: {domainName: query}})
        .then((response) => console.log(response))
        .catch((error) => console.log(error))
    }
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
}

#app img {
  width: 200px;
}

#app input.input-error {
  border-color: red;
}
</style>
