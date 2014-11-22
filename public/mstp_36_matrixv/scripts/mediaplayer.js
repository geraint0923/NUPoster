function matrix_media_player(){
		
    jQuery("#jquery_jplayer_1").jPlayer({
        ready: function () {
            jQuery(this).jPlayer("setMedia", {
                oga : "source of file",
				mp3 : "source of file"
            });
        },
        swfPath: "/scripts",
		cssSelectorAncestor: "#jp_container_1",
        supplied: "oga,mp3",
    });
	
	//To create more videos, just duplicate the codes above. However, remember to change #jquery_jplayer_1 and #jp_container_1 to something else which matches the ID set at the tile preview / single pages.


}