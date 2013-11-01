Messages = new Meteor.Collection('messages');

if (Meteor.isClient) {
  Template.login.events({
    'keydown input#user_name' : function(ev) {
      if (ev.keyCode == 13) {
        Session.set('user_name', ev.target.value);
      }
    }
  });

  Handlebars.registerHelper('user_name', function() {
    return Session.get('user_name');
  });

  Template.chat.messages = function() {
    return Messages.find({}, {sort: {timestamp: 1}});
  }

  Handlebars.registerHelper('time', function(timestamp) {
    var date = new Date(timestamp);
    var pad = function(n) { return n < 10 ? "0" + n : n.toString() };
    return pad(date.getHours()) + ":" + pad(date.getMinutes());
  });

  Template.chat.events({
    'keydown input#message' : function(ev) {
      if (ev.keyCode == 13 && ev.target.value.length > 0) {
        Messages.insert({name: Session.get('user_name'),
                         text: ev.target.value,
                         timestamp: new Date().getTime()});
        ev.target.value = "";
      }
    }
  });
}
