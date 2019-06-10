<template>
  <div id="app">
    <img class="app-logo" alt="Vue logo" src="./assets/logo.png">
    <HelloWorld msg="Welcome to Server ID Service"/>
    <SearchBox
      plh="Type a server domain"
      v-on:search="searchServer"
    />
    <DomainCard v-if="currentDomain" v-bind:domain="currentDomain"/>
    <DomainList v-bind:domains="domains"/>
  </div>
</template>

<script>
// Custom Components
import HelloWorld from './components/HelloWorld.vue'
import SearchBox  from './components/SearchBox.vue'
import DomainCard from './components/DomainCard.vue'
import DomainList from './components/DomainList.vue'

// Libraries
import axios      from 'axios'

export default {
  name: 'app',
  data: function() {
    return {
      currentDomain: null,
      domains: {}
    }
  },
  components: {
    HelloWorld,
    SearchBox,
    DomainCard,
    DomainList
  },
  created: function() {
    this.fetchDomains();
  },
  methods: {
    searchServer: function(query) {
      axios.get("http://localhost:3000/domains/search", {params: {domainName: query}})
        .then((response) => this.currentDomain = response.data )
        .then(this.fetchDomains)
        .catch((error) => console.log(error))
    },
    fetchDomains: function() {
      axios.get("http://localhost:3000/domains")
        .then((response) => this.domains = response.data)
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

#app img.app-logo {
  width: 200px;
}

#app input.input-error {
  border-color: red;
}
</style>
