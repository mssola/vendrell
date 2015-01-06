/*
 * Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
 * This file is licensed under the MIT license.
 * See the LICENSE file.
 */


/*
 * On load setup the event listener of the "Remove player" dialog. If the user
 * has typed the name of the user, then he is allowed to remove this player.
 */
window.onload = function() {
    var name = document.getElementById('name').value;
    var el = document.getElementById('rm-text');
    var btn = document.getElementById('rm-btn');

    el.onkeyup = function() { btn.disabled = !(el.value === name); }
}

