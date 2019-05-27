function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

function start(e) {
  e.preventDefault();

  // Normalize the various vendor prefixed versions of getUserMedia.
  navigator.getUserMedia = (navigator.getUserMedia ||
                            navigator.webkitGetUserMedia ||
                            navigator.mozGetUserMedia || 
                            navigator.msGetUserMedia);

	// Check that the browser supports getUserMedia.
	// If it doesn't show an alert, otherwise continue.
	if (navigator.getUserMedia) {
	  // Request the camera.
	  navigator.getUserMedia(
	    // Constraints
	    {
	      audio: true,
	      video: true
	    },

	    // Success Callback
	    function(localMediaStream) {
		// Get a reference to the video element on the page.
		var video = document.getElementById('camera-stream');

		// Create an object URL for the video stream and use this 
		// to set the video source.
		video.srcObject = localMediaStream

		// returns a frame encoded in base64
		const getFrame = () => {
		    const canvas = document.createElement('canvas');
		    canvas.width = video.videoWidth;
		    canvas.height = video.videoHeight;
		    canvas.getContext('2d').drawImage(video, 0, 0);
		    const data = canvas.toDataURL('image/png');
		    return data;
		}

		var id = uuid()
		var proto = (window.location.protocol == "http:") ? "ws://" : "wss://"
		var url = proto + window.location.host + window.location.pathname + "/video?id=" + id;
		var share = document.getElementById('share');

		const FPS = 9;
		const ws = new WebSocket(url);
		var send;

		ws.onopen = () => {
		    console.log(`Connected to ${ws}`);
		    send = setInterval(() => {
		        ws.send(getFrame());
		    }, 1000 / FPS);
		}

		const stop = (e) => {
  		  e.preventDefault();
		  var stream = video.srcObject;
		  var tracks = stream.getTracks();

		  for (var i = 0; i < tracks.length; i++) {
		    var track = tracks[i];
		    track.stop();
		  }

		  video.srcObject = null;
		  clearInterval(send)
		  ws.close();
		  share.innerHTML = '';
		}

		// stop stream
		document.getElementById('stop').onclick = stop;

		// set share link
		shareURL = window.location.href + "?id=" + id + "&type=client"
		share.innerHTML = 'Share your stream <a href="'+ shareURL + '">' + shareURL + '</a>';
	    },

   
	    // Error Callback
	    function(err) {
	      // Log the error to the console.
	      console.log('The following error occurred when trying to use getUserMedia: ' + err);
	    }
	  );

	} else {
	  alert('Sorry, your browser does not support getUserMedia');
	}

}

function getParam(name) {
    name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]');
    var regex = new RegExp('[\\?&]' + name + '=([^&#]*)');
    var results = regex.exec(location.search);
    return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
};

window.onload = function() {
	var front = false;

	var id = getParam("id");
	var type = getParam("type");

	if (id.length > 0 && type == "client") {
		document.getElementById('control').innerHTML = '';

		var vid = document.getElementById('video-container');
		var vcc = document.getElementById('camera-stream');
		var img = document.createElement('img');
		vid.removeChild(vcc);
		vid.prepend(img);
		
		var proto = (window.location.protocol == "http:") ? "ws://" : "wss://"
		var url = proto + window.location.host + window.location.pathname + "/video" + window.location.search;

		// new websocket
		const ws = new WebSocket(url);

		ws.onopen = () => console.log(`Connected to ${url}`);
		ws.onmessage = message => {
		    // set the base64 string to the src tag of the image
		    img.src = message.data;
		}
	} else {
		// start stream
		document.getElementById('start').onclick = start;
	}

	// flip camera
	//document.getElementById('flip').onclick = function() { front = !front; };
	//var constraints = { video: { facingMode: (front? "user" : "environment") } };
}
