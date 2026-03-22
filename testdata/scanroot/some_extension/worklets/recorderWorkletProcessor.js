/*
A worklet for recording in sync with AudioContext.currentTime.

More info about the API:
https://developers.google.com/web/updates/2017/12/audio-worklet

Based on @flpvsk implementation, addapted to use with multi-channeled audio streams.

For a full example of my use look here: 
https://github.com/theroman/Sample/blob/master/js/external/recorder/recorder.js


How to use:

1. Serve this file from your server (under "public/worklets") as is.

2. Register the worklet:

    const audioContext = new AudioContext();
    await audioContext.audioWorklet.addModule('path/to/recorderWorkletProcessor.js')

3. Prepare to record:

    # Create worklet

    const channels = 2
    const recorderNode = new window.AudioWorkletNode(audioContext, 
                                                     'recorder-worklet', 
                                                     {parameterData: {numberOfChannels: channels}});

    # Connect your source

    yourSourceNode.connect(recorderNode);
    recorderNode.connect(audioContext.destination);

    # Register worklet events
    
    recorderNode.port.onmessage = (e) => {
        const data = e.data;
        switch(data.eventType) {
          case "data":
            // process pcm data; encode etc
            const audioData = data.audioBuffer;
            const bufferSize = data.bufferSize;
            break;
          case "stop":
            // recording has stopped
            break;
        }
    };

4. Start the recording 
    let isRecording = recorderNode.parameters.get('isRecording')
    isRecording.setValueAtTime(1, audioContext.currentTime);

5. Stop the recording
    let isRecording = recorderNode.parameters.get('isRecording')
    isRecording.setValueAtTime(0, audioContext.currentTime);
      
*/

class RecorderWorkletProcessor extends AudioWorkletProcessor {
    static get parameterDescriptors() {
      return [{
        name: 'isRecording',
        defaultValue: 0
      },
      {
        name: 'numberOfChannels',
        defaultValue: 2
      },
      {
        name: 'channelsToRecord',
        defaultValue: 2
      }
    ];
    }
  
    constructor() {
      super();
      this._bufferSize = 2048;
      this._buffers = null;
      this._initBuffer();
      this.total_flush = 0
    }

    _initBuffers(numberOfChannels) {
      this._buffers = []
      for (let channel=0; channel < numberOfChannels; channel++) {
        this._buffers.push(new Float32Array(this._bufferSize))
      }
    }
  
    _initBuffer() {
      this._bytesWritten = 0;
    }
  
    _isBufferEmpty() {
      return this._bytesWritten === 0;
    }
  
    _isBufferFull() {
      return this._bytesWritten === this._bufferSize;
    }


    _pushToBuffers(audioRawData, channelsToRecord) {
      if (this._isBufferFull()) {
          this._flush();
      }

      let dataLength = audioRawData[0].length
      for (let idx=0; idx<dataLength; idx++) {
        let leftValue = audioRawData[0][idx]
        let rightValue = audioRawData[1][idx]
        if (channelsToRecord == 2) { // both
          this._buffers[0][this._bytesWritten] = leftValue
          this._buffers[1][this._bytesWritten] = rightValue
        }
        if (channelsToRecord == 0) { // left
          this._buffers[0][this._bytesWritten] = leftValue
          this._buffers[1][this._bytesWritten] = leftValue
        }
        if (channelsToRecord == 1) { // right
          this._buffers[0][this._bytesWritten] = rightValue
          this._buffers[1][this._bytesWritten] = rightValue
        }
        this._bytesWritten += 1
      }
    }
  
    _flush() {
      let buffers = []
      this._buffers.forEach((buffer, channel) => {
        if (this._bytesWritten < this._bufferSize) {
          buffer = buffer.slice(0, this._bytesWritten);
        }
        buffers[channel] = buffer
      })
      this.port.postMessage({
        eventType: 'data',
        audioBuffer: buffers,
        bufferSize: this._bufferSize
      });
      this.total_flush += 1
      this._initBuffer();
    }
  
    _recordingStopped() {
      this.port.postMessage({
        eventType: 'stop'
      });
    }
  
    process(inputs, outputs, parameters) {
      const isRecordingValues = parameters.isRecording;
      const numberOfChannels = parameters.numberOfChannels[0]   
      const channelsToRecord = parameters.channelsToRecord[0]
      if (this._buffers === null) {
        this._initBuffers(numberOfChannels)
      }
      
      for (let dataIndex = 0; dataIndex < isRecordingValues.length; dataIndex++) 
      {
        const shouldRecord = isRecordingValues[dataIndex] === 1;
        if (!shouldRecord && !this._isBufferEmpty()) {
          this._flush();
          this._recordingStopped();
        }
  
        if (shouldRecord) {
          let audioRawData = inputs[0]
          this._pushToBuffers(audioRawData, channelsToRecord)
        }
      }
      return true;
    }
  
  }
  
  registerProcessor('recorder-worklet', RecorderWorkletProcessor);