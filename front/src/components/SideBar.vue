<template>
  <div class="sidebar">
    <img class="app-logo" alt="Vue logo" src="../assets/logo.png">
    <ul>
      <li class="active">
        <svg class="icon">
          <use xlink:href="../assets/icons.svg#server"/>
      </svg>
        Servers
      </li>
      <li>
        <svg class="icon">
          <use xlink:href="../assets/icons.svg#server"/>
      </svg>
        About
      </li>
    </ul>
  </div>
</template>

<script>

export default {
  name: 'HelloWorld',
  props: {
    msg: String
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.sidebar {
  padding: 20px 10px;
  text-align: center;
  box-shadow: inset -5px 0px 4px 0px #e5e9ec;
}

ul {
  list-style-type: none;
  padding: 0;
}

.sidebar img.app-logo {
  margin-bottom: 50px;
  width: 150px;
}

.sidebar li {
  margin: 0 auto 20px;
  position: relative;
  padding: 10px 0;
  cursor: pointer;
  display: block;
  color: #e8eaec;
}

.sidebar li.active {
  color: #8b97a3;
}

.sidebar li.active:after {
  background: linear-gradient(#5bc15a, #49bbac);
  position: absolute;
  display: block;
  content: "";
  height: 100%;
  width: 5px;
  top: 0;
  right: -10px;
}

.sidebar svg.icon {
  margin: 0 auto 10px;
  display: block;
  fill: #e8eaec;
  height: 40px;
  width: 40px;
}

.sidebar li.active svg.icon {
  fill: #8b97a3;
}

</style>
