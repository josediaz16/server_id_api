<template>
  <div class="sidebar">
    <img class="app-logo" alt="Vue logo" src="../assets/logo.png">
    <ul>
      <li class="active">
        <Icon iconname="server"/>
        Servers
      </li>
      <li>
        <Icon iconname="improvement"/>
        Comming soon
      </li>
    </ul>
  </div>
</template>

<script>

import Icon from './Icon'

export default {
  name: 'SideBar',
  props: {
    msg: String
  },
  components: {
    Icon
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss">

.sidebar {
  box-shadow: inset -5px 0px 4px 0px #e5e9ec;
  padding: 20px 10px;
  text-align: center;

  img.app-logo {
    margin-bottom: 50px;
    width: 130px;
  }

  li {
    margin: 0 auto 20px;
    position: relative;
    padding: 10px 0;
    cursor: pointer;
    display: block;
    color: #e8eaec;

    &.active {
      color: #8b97a3;

      &:after {
        background: linear-gradient(#5bc15a, #49bbac);
        position: absolute;
        display: block;
        content: "";
        height: 100%;
        width: 5px;
        top: 0;
        right: -10px;
      }

      svg.icon {
        fill: #8b97a3;
      }
    }
  }

  svg.icon {
    margin: 0 auto 10px;
    display: block;
    fill: #e8eaec;
    height: 40px;
    width: 40px;
  }
}

</style>
