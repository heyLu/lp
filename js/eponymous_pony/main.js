Names = new Meteor.Collection("names");

if (Meteor.isClient) {
  Template.names.helpers({
    names: function() {
      return Names.find({}).fetch();
    }
  });

  Template.names.events({
    'keypress input': function(event) {
      if (event.key == "Enter") {
        Names.insert({name: event.currentTarget.value});
        event.currentTarget.value = "";
      }
    }
  });
}