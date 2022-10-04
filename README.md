#### Превьювер изображений
Данный сервис - web-сервер, загружающий изображения и масштабирующий/обрезающий их до нужного формата. 
Наиболее частые изображения остаются в кэше, редко заменяются новыми.

#### Поддерживаемые форматы
Поддерживаются форматы файлов: jpeg, png, gif

### Параметры URL запроса
http://address:port/service/width/height/somesite.com/image.jpg

Параметры:
address:port - адрес:порт
service - тип операции над изображением
- "resize" - изменить размер, пропорции не сохраняются.
- "fit" - вписать изображение в заданный размер, пропорции сохраняются,
- "fill" - заполнить, пропорции сохраняются,
исходное изображение при этом центрируется и может быть обрезано по высоте или ширине
Полученное изображение может быть меньше по высоте или ширине
width - ширина
height - высота
site.com/image.jpg - адрес до изображения на стороннем ресурсе
Выбор тип протокола (https, http), происходит автоматически. Сначала проверяет url по https, затем в случае ошибки по http.

#### Пример Url для масштабирования изображения
http://0.0.0.0:8013/resize/80/80/www.yougiveme.ru/web/upload/images/86/23183/offer_92_1664807898.jpeg

#### Компиляция, запуск, тестирование
make - скомпилировать проект, выходная папка ./bin
make run - собрать и запустить докер образ
make test - запустить юнит и интеграционные тесты

#### Конфигурация
Настройка осуществляется на основании ENV переменных окружения, с предопределёнными дефолтными значениями:

LOG_LEVEL - уровень логирования ("error", "warn", "info", "debug"), по умолчанию: "debug"
HTTP_LISTEN - адрес:порт, на котором запущен сервис, по умолчанию: ":8013"
IMAGE_MAX_FILE_SIZE - максимальный размер запрашиваемого (исходного) изображения в байтах,
по умолчанию: "1000000" (1M)
IMAGE_GET_TIMEOUT - максимальное время в секундах, в течение которого сервис будет пытаться получить
удалённое изображение, по умолчанию: 10 сек
CACHE_SIZE - общий размер кэша для всех обработанных и изменённых картинок в байтах по умолчанию: "100000000" (100M)
CACHE_TYPE - тип кэша ("inmemory" - в оперативной памяти, "disk" - указанная папка на диске), по умолчанию: "disk"
CACHE_PATH - путь к папке кэша на диске, по умолчанию "./cache"