<template>
  <div class="domain-card">
    <img v-bind:src="domain.logo" alt="">
    <h4>{{domain.title}}</h4>
    <ul>
      <li>Domain Up: {{!domain.is_down}}</li>
      <li>Ssl Grade: {{domain.ssl_grade}}</li>
      <li>Servers Changed: {{domain.servers_changed}}</li>
      <li>Previous Ssl Grade: {{previousSslGrade}}</li>
    </ul>

    <h5>Servers</h5>
    <div v-for="(server, index) in domain.servers" v-bind:key="server.ip_address" class="server-item">
      <h6>{{index}}</h6>
      <ul>
        <li>Address: {{server.ip_address}}</li>
        <li>Ssl Grade: {{server.grade}}</li>
        <li>Country: {{server.country}}</li>
        <li>Owner: {{server.owner}}</li>
      </ul>
    </div>
  </div>
</template>

<script>

export default {
  name: "DomainCard",
  props: {
    domain: Object
  },
  computed: {
    previousSslGrade: function() {
      return this.domain.previous_ssl_grade === "" ? "Unknown" : this.domain.previous_ssl_grade
    }
  }
}
</script>

<style>
  .domain-card img {
    width: 35px;
  }
</style>
