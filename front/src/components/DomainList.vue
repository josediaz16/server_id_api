<template>
  <div class="domain-list">
    <div
      class="domain-item"
      v-for="(domain, name) in domains"
      v-bind:key="name"
    >
      <h4>{{name}}</h4>
      <DomainCard v-bind:domain="domain"/>
    </div>
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
