$(function () {

  'use strict';

  $('.qor-locales').on('change', function () {
    window.location.assign($(this).val());
  });

});
