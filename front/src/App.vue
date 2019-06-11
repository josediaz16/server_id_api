<template>
  <div id="app">
    <SideBar />
    <div class="dashboard">
      <SearchBox
        plh="Type a server domain"
        v-on:search="searchServer"
      />
      <DomainCard v-if="currentDomain" v-bind:domain="currentDomain"/>
      <DomainList v-bind:domains="domains"/>
    </div>
  </div>
</template>

<script>
// Custom Components
import SideBar    from './components/SideBar.vue'
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
    SideBar,
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

body {
  margin: 0;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;

  display: grid;
  grid-template-columns: 1fr 9fr;
}

#app .dashboard {
  background: #eaf0f5;
  padding: 20px;
}

#app input.input-error {
  border-color: red;
}
</style>
