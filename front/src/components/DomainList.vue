<template>
  <div class="domain-list">
    <DomainCard
      v-for="(domain, name) in domains"
      v-bind:key="name"
      v-bind:domain="domain"
      v-bind:name="name"
    />
  </div>
</template>

<script>

import DomainCard from './DomainCard.vue'

export default {
  name: 'DomainList',
  props: {
    domains: Object
  },
  components: {
    DomainCard
  }
}
</script>

<style>
.domain-list {
  grid-template-columns: 1fr 1fr;
  grid-gap: 30px;
  display: grid;
}

</style>
