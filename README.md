# Запуск

## kubectl get pods
```
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
goapp-7f7c48cb4-6x2mk       1/1     Running   0          29s
goapp-7f7c48cb4-lhb2c       1/1     Running   0          43s
goapp-7f7c48cb4-lskgk       1/1     Running   0          56s
postgres-7775db7ddc-28qd2   1/1     Running   0          5m38s
redis-5d85cc97f8-jrpbc      1/1     Running   0          5m36s
```

## kubectl get services
```
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> kubectl get services
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
goapp        ClusterIP   10.103.22.233   <none>        80/TCP     7m51s
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP    8m43s
postgres     ClusterIP   10.110.202.24   <none>        5432/TCP   7m54s
redis        ClusterIP   10.106.20.16    <none>        6379/TCP   7m52s
```

## curls после проброски порта
```
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> curl.exe http://localhost:8080/healthz
{"status":"ok"}
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> curl.exe http://localhost:8080/readyz
{"status":"ready"}
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> curl.exe http://localhost:8080/users
[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"},{"id":3,"name":"Charlie"}]
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> curl.exe "http://localhost:8080/cache/set?key=hello&value=world"
{"key":"hello","status":"ok","value":"world"}
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> curl.exe "http://localhost:8080/cache/get?key=hello"
{"key":"hello","value":"world"}
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4>
```

## logs 
```
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> kubectl logs goapp-7f7c48cb4-6x2mk --tail=20
{"time":"2026-10-01T16:59:16.871692602Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":28602,"remote":"10.244.0.1:37070"}
{"time":"2026-10-01T16:59:17.245535995Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":500524,"remote":"10.244.0.1:37078"}
{"time":"2026-10-01T16:59:22.24521669Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":479623,"remote":"10.244.0.1:50010"}
{"time":"2026-10-01T16:59:26.871514221Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":27901,"remote":"10.244.0.1:50016"}
{"time":"2026-10-01T16:59:27.245036905Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":499424,"remote":"10.244.0.1:50026"}
{"time":"2026-10-01T16:59:32.244631373Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":489823,"remote":"10.244.0.1:47818"}
{"time":"2026-10-01T16:59:36.872004357Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":31102,"remote":"10.244.0.1:47834"}
{"time":"2026-10-01T16:59:37.243734629Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":482822,"remote":"10.244.0.1:47848"}
{"time":"2026-10-01T16:59:42.244614519Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":527925,"remote":"10.244.0.1:47316"}
{"time":"2026-10-01T16:59:46.870826708Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":25501,"remote":"10.244.0.1:47326"}
{"time":"2026-10-01T16:59:47.244708052Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":509124,"remote":"10.244.0.1:47332"}
{"time":"2026-10-01T16:59:52.244924148Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":712734,"remote":"10.244.0.1:59118"}
{"time":"2026-10-01T16:59:56.871489355Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":28901,"remote":"10.244.0.1:59124"}
{"time":"2026-10-01T16:59:57.244583506Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":574026,"remote":"10.244.0.1:59140"}
{"time":"2026-10-01T17:00:02.243015584Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":545324,"remote":"10.244.0.1:51702"}
{"time":"2026-10-01T17:00:06.869291495Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":24501,"remote":"10.244.0.1:51716"}
{"time":"2026-10-01T17:00:07.243094602Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":576428,"remote":"10.244.0.1:51730"}
{"time":"2026-10-01T17:00:12.242520289Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":487623,"remote":"10.244.0.1:44216"}
{"time":"2026-10-01T17:00:16.87005584Z","level":"INFO","msg":"http request","method":"GET","path":"/healthz","status":200,"duration":26601,"remote":"10.244.0.1:44220"}
{"time":"2026-10-01T17:00:17.243293947Z","level":"INFO","msg":"http request","method":"GET","path":"/readyz","status":200,"duration":502024,"remote":"10.244.0.1:44236"}
```

## Масштабирование
```
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> kubectl scale deployment goapp --replicas=5
deployment.apps/goapp scaled
```

## Обновление образа

Тут myregistry/goapp:v2.0 не существует и по идее результат будет не тот, который ожидается. Поэтому я здесь соберу реальный v2 образ, если надо по-другому - переделаю.

Выполню docker build -t goapp:v2 .

После minikube image load goapp:v2

```
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> minikube image load goapp:v2
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> kubectl set image deployment/goapp server=goapp:v2
deployment.apps/goapp image updated
PS C:\Users\maksn\OneDrive\Desktop\Папка с папками\Работа\EM Tasks\6. DevOps\Task4> kubectl rollout status deployment/goapp
deployment "goapp" successfully rolled out
```