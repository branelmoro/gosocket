package gosocket

type onBadRequest func(HttpRequest)

type onRequest func(HttpWriter, HttpRequest)

type wsWriterCb func(WsWriter)

type messageCb func(WsWriter, Message)

type textCb func(WsWriter, string)

type binaryCb func(WsWriter, []byte)

type closeCb func(WsWriter, CloseMsg)

type errorCb func(WsWriter, error)

var(
    OnMalformedRequest onBadRequest
    OnHttpRequest onRequest
    OnWebsocketRequest onRequest
    OnWebsocketOpen wsWriterCb
    OnMessage messageCb

    OnText textCb
    OnBinary binaryCb

    OnClose closeCb
    OnError errorCb
    OnPing wsWriterCb
    OnPong wsWriterCb

    // OnText messageCb
    // OnBinary messageCb

)



// OnMultiFrameMessage onFrame

// for {
// 	if frame.Size() > 23424324 {
// 		frameData, err = frame.FetchData()
// 		frame.WriteTo
// 		if err != nil {
// 			// fetch data error
// 			w.Close()
// 			return
// 		}
// 	} else {
// 		// frame size error
// 		w.Close()
// 		return
// 	}
// 	if frame.Final() {
// 		break
// 	}
// 	frame, err = frame.Next()
// 	if err != nil {
// 		// next frame error
// 		w.Close()
// 		return
// 	}
// }