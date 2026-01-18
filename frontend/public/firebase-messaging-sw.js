importScripts('https://www.gstatic.com/firebasejs/10.7.1/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.7.1/firebase-messaging-compat.js');

firebase.initializeApp({
  apiKey: "AIzaSyCP-CjGEcfqigAWLIHwC7q8tYVHKAGtL0w",
  authDomain: "guss-sns.firebaseapp.com",
  projectId: "guss-sns",
  storageBucket: "guss-sns.firebasestorage.app",
  messagingSenderId: "544714373795",
  appId: "1:544714373795:web:60a555df7308cc82265d16"
});

const messaging = firebase.messaging();

// 백그라운드 메시지 수신
messaging.onBackgroundMessage((payload) => {
  console.log('[firebase-messaging-sw.js] 백그라운드 메시지 수신:', payload);
  
  const notificationTitle = payload.notification.title;
  const notificationOptions = {
    body: payload.notification.body,
    icon: '/favicon.ico'
  };

  self.registration.showNotification(notificationTitle, notificationOptions);
});