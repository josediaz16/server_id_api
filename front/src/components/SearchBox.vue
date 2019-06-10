<template>
  <div class="search-box">
    <input v-model="searchText" type="search" :placeholder="plh"/>
    <button v-on:click="searchServer">Search</button>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'SearchBox',
  data: () => {
    return {
      searchText: ''
    }
  },
  props: {
    plh: String,
  },
  methods: {
    searchServer: function() {
      axios.get("http://localhost:3000/domains/search", {params: {domainName: this.searchText}})
        .then((response) => console.log(response))
        .catch((error) => console.log(error))
    }
  }
}
</script>
