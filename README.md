# README [![Build Status](https://travis-ci.org/mssola/vendrell.svg?branch=master)](https://travis-ci.org/mssola/vendrell)

This is a web application that I created for a friend of mine. This web app helps
him to manage a group of hockey players. This application provides the
following clean and simple work flow:

The admin user can create/remove/rename players. Each player has a page
(corresponds to the "show" action in a REST architecture). When a player
accesses to his page, it can rate the last practice. This is basically all that
a player can do in this application.

The admin can fetch the ratings from two different pages: the root page and a
player page. The root page has all the ratings from all the players on the
system. If the admin accesses a player page, then he fetches the ratings from
this specific player.

Last but not least, the admin user can download a CSV file containing ratings
from either the root page or a player page.

Copyright &copy; 2014 Miquel Sabaté Solà, released under the MIT License.

