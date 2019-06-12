<template>
  <div id="app">
    <SideBar />
    <div class="dashboard">
      <DomainList v-bind:domains="domains"/>
    </div>
    <div>
      <SearchBox
        plh="Type a server domain"
        v-on:search="searchServer"
      />
      <DomainCard v-if="currentDomain" v-bind:domain="currentDomain"/>
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
        .catch(console.log)
    },
    fetchDomains: function() {
      axios.get("http://localhost:3000/domains")
        .then((response) => this.domains = response.data)
        .catch(console.log)
    }
  }
}
</script>

<style>

body {
  margin: 0;
}

h4, h5 {
  margin: 10px 0;
  color: #5d6777;
}

h5 {
  font-size: 0.95em;
}

ul {
  list-style-type: none;
  padding: 0;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;

  display: grid;
  grid-template-columns: 150px 1fr 350px;
}

#app .dashboard {
  background: #eaf0f5;
  padding: 20px 50px;
}

#app input.input-error {
  border-color: red;
}
</style>
