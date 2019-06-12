<template>
  <div class="search-box">
    <input
      v-model.trim="searchText"
      type="search"
      v-on:input="validateText"
      :placeholder="plh"
      v-bind:class="inputClass"
    />
    <button v-on:click="performSearch">Search</button>
    <span v-if="touched && !validDomain">Please check the domain</span>
  </div>
</template>

<script>

const domainRegex = /[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+/

export default {
  name: 'SearchBox',
  data: () => {
    return {
      searchText: '',
      queryText: '',
      validDomain: false,
      touched: false
    }
  },
  props: {
    plh: String,
  },
  methods: {
    validateText: function() {
      const match = this.searchText.match(domainRegex)

      if (match) {
        this.queryText = match[0]
        this.validDomain = true
      }
      else {
        this.validDomain = false
      }
    },
    performSearch: function() {
      this.touched = true

      if (this.validDomain) {
        this.$emit('search', this.queryText)
      }
    }
  },
  computed: {
    inputClass: function() {
      return {
        'input-error': this.touched && !this.validDomain
      }
    }
  }
}
</script>

<style lang="scss">
.search-box {
  grid-template-columns: 5fr 1fr;
  grid-column-gap: 10px;
  padding: 20px 5px;
  display: grid;

  input {
    background: #f6f8fa;
    border-radius: 5px;
    padding: 10px;
    border: none;

    &.input-error {
      border: 1px solid #d75b4b;
    }
  }

  button {
    background: rgba(79, 179, 103, 0.73);
    border-radius: 5px;
    font-weight: bold;
    color: #fff;
    border: none;
  }

  input.input-error ~ span {
    font-size: 13px;
    color: #d75b4b;
    padding: 5px 0;
  }
}

</style>
