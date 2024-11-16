
INSERT INTO users (name, email, login, password, is_admin) VALUES
('Timofei Tsypyshev', 'timofeitsypyshev@yandex.ru', 'admin1', 'password123', TRUE),
(NULL, NULL, 'user1', 'userPass123', FALSE),
(NULL, NULL, 'user2', 'securePass456', FALSE),
(NULL, NULL, 'user3', 'myPassword789', FALSE);



INSERT INTO langs (name, img_link, short_description, author, year, version, description, list, status) VALUES
('Python', 'http://localhost:9000/code-inspector/python.png', 
 'Объединяет простоту и мощь', 'Гвидо ван Россум', '1991', 
 'Python 3.12.6 (Sep 12, 2024)', 
 '— это высокоуровневый язык программирования общего назначения, который широко используется благодаря своей гибкости, простоте и мощным возможностям расширения. Вот основные технические характеристики Python:', 
 '{"Исполняемость": "Python интерпретируемый язык. Код выполняется интерпретатором, который читает и исполняет команды строка за строкой.",
  "Модульность и библиотеки": "Python поддерживает создание модулей и пакетов, а также имеет обширную стандартную библиотеку с множеством функциональных возможностей.",
  "Мультисистемность": "Python доступен для большинства операционных систем, включая Windows, macOS и Linux.",
  "Объектно-Ориентированное Программирование (ООП)": "Python поддерживает ООП, позволяя создавать классы и объекты.",
  "Пакетный менеджер": "Python использует pip в качестве стандартного инструмента для установки и управления пакетами и библиотеками.",
  "Производительность": "Python может быть медленнее по сравнению с компилируемыми языками, но его производительность можно улучшить с помощью различных оптимизаций, таких как Cython и PyPy.",
  "Синтаксис": "Python имеет простой и читаемый синтаксис, что делает его удобным для новичков. Отступы используются для обозначения блоков кода.",
  "Типизация": "Python является динамически типизированным языком, что означает, что типы переменных проверяются во время выполнения программы."}', TRUE),

('C++', 'http://localhost:9000/code-inspector/cpp.png', 
 'Контроль и производительность в одном лице', 'Бьёрн Страуструп', '1985', 
 'C++20 (Dec 2020)', 
 '— это мощный и высокопроизводительный язык программирования, известный своей способностью обеспечивать низкоуровневый доступ к памяти и поддерживать сложные структуры данных. Вот основные технические характеристики C++:', 
 '{"Производительность": "C++ обеспечивает высокую производительность за счет компиляции в машинный код и эффективного управления ресурсами.",
  "Синтаксис": "C++ поддерживает широкий спектр программных конструкций, включая функции, классы, шаблоны и перегрузку операторов, что делает его подходящим для различных типов программирования.",
  "Объектно-Ориентированное Программирование (ООП)": "C++ позволяет создавать и использовать классы и объекты, поддерживает наследование, полиморфизм и инкапсуляцию.",
  "Шаблоны": "C++ поддерживает шаблоны для создания обобщенного кода, который может работать с различными типами данных.",
  "Мультисистемность": "C++ доступен для множества операционных систем и платформ, что позволяет создавать кроссплатформенные приложения.",
  "Стандартная библиотека": "C++ включает стандартную библиотеку, которая предоставляет функции и классы для работы с коллекциями, вводом/выводом, и алгоритмами.",
  "Управление памятью": "C++ позволяет явное управление памятью с помощью указателей и динамического выделения, что требует от программиста внимательного контроля ресурсов.",
  "Проектирование и сборка": "C++ проектирование и сборка поддерживаются различными системами сборки и управления зависимостями, такими как CMake и Make."}', TRUE),

('GO', 'http://localhost:9000/code-inspector/golang.png', 
 'Эффективный для масштабируемых решений', 'Роберт Гризмер, Роб Пайк', '2009', 
 'Go 1.21.0 (Aug 2023)', 
 '— это язык программирования, разработанный Google с акцентом на простоту, производительность и эффективное параллельное выполнение. Вот основные технические характеристики Go:', 
 '{"Типизация": "Go является статически типизированным языком, что означает проверку типов переменных на этапе компиляции.",
  "Исполняемость": "Go компилируемый язык. Код Go компилируется в машинный код, что обеспечивает высокую производительность выполнения программ.",
  "Синтаксис": "Go имеет простой и лаконичный синтаксис, который делает его легким для изучения и использования. Отступы используются для форматирования кода, что способствует чистоте и читаемости кода.",
  "Параллелизм": "Go предоставляет встроенную поддержку для параллельного выполнения через горутины и каналы, что упрощает разработку многопоточных приложений.",
  "Мультисистемность": "Go поддерживает кроссплатформенность и доступен для большинства операционных систем, включая Windows, macOS и Linux.",
  "Модульность и библиотеки": "Go имеет встроенную поддержку для работы с пакетами и модулями, а также поставляется с обширной стандартной библиотекой.",
  "Пакетный менеджер": "Go использует собственный инструмент go для управления пакетами и зависимостями, который упрощает процесс сборки и установки.",
  "Производительность": "Go обеспечивает высокую производительность благодаря компиляции в машинный код и оптимизированному управлению памятью, что делает его эффективным для сетевых и серверных приложений."}', TRUE),

('HTML', 'http://localhost:9000/code-inspector/html.png', 
 'Основа структуры и содержания веб-страниц', 'Тим Бернерс-Ли', '1993', 
 'HTML5 (Oct 2014)', 
 '— это стандартный язык разметки, используемый для создания и структурирования веб-страниц. Вот основные технические характеристики HTML:', 
 '{"Разметка": "HTML использует теги для определения структуры и содержания веб-страниц. Теги обозначают различные элементы, такие как заголовки, параграфы, ссылки и изображения.",
  "Исполняемость": "HTML не является исполняемым языком программирования; вместо этого он интерпретируется веб-браузерами, которые отображают содержимое страницы согласно разметке.",
  "Синтаксис": "HTML имеет простой и гибкий синтаксис, который легко читается и пишется. Теги обычно имеют открывающую и закрывающую форму, например <p> и </p> для параграфов.",
  "Структура": "HTML обеспечивает основу для веб-страниц, включая такие элементы, как заголовок, тело и разделение на секции. Он также поддерживает вложенность элементов для создания сложных структур.",
  "Мультисистемность": "HTML является стандартом для веб-разработки и поддерживается всеми современными веб-браузерами, независимо от операционной системы.",
  "Модульность и расширения": "HTML можно расширять с помощью CSS (Cascading Style Sheets) для стилизации и JavaScript для интерактивности и динамического поведения.",
  "Пакетный менеджер": "HTML сам по себе не имеет пакетного менеджера, но может быть использован вместе с инструментами и фреймворками для упрощения веб-разработки, такими как npm (Node Package Manager) для JavaScript-библиотек.",
  "Производительность": "Поскольку HTML представляет собой статический язык разметки, его производительность в основном зависит от браузера и качества кода, используемого для поддержки стилей и сценариев."}', TRUE),

('CSS', 'http://localhost:9000/code-inspector/css.png', 
 'Создает оформление веб-интерфейсов', 'Грабриел Маззоне', '1996', 
 'CSS3 (2011)', 
 '— это язык стилей, используемый для управления внешним видом и форматированием веб-страниц, созданных на HTML. Вот основные технические характеристики CSS:', 
 '{"Стилевое оформление": "CSS применяется для стилизации HTML-элементов, позволяя задавать цвета, шрифты, отступы, размеры и другие визуальные свойства.",
  "Синтаксис": "CSS имеет простой синтаксис, состоящий из селекторов и деклараций. Селекторы указывают элементы, к которым применяются стили, а декларации определяют свойства и их значения.",
  "Каскадность и наследование": "CSS поддерживает каскадность, что означает, что стили могут наследоваться от родительских элементов и могут переопределяться более специфичными правилами.",
  "Респонсивный дизайн": "CSS предоставляет возможности для создания адаптивных веб-дизайнов, используя медиа-запросы и флексбоксы для управления расположением элементов на экране в зависимости от размера устройства.",
  "Мультисистемность": "CSS поддерживается всеми современными веб-браузерами и может использоваться с HTML для создания интерактивных и визуально привлекательных веб-страниц.",
  "Расширяемость": "CSS может быть расширен с помощью препроцессоров, таких как SASS и LESS, которые добавляют дополнительные возможности и упрощают организацию стилей.",
  "Производительность": "Поскольку CSS управляет визуальными элементами на веб-странице, его производительность в основном зависит от браузера и качества кода, который используется для стилизации."}', TRUE),

('COBOL', NULL, 
 'Язык для бизнес-приложений', 'Грейс Хоппер', '1959', 
 'COBOL 2002', 
 '— это язык программирования, разработанный для обработки бизнес-данных. Он был одним из первых языков, специально созданных для бизнес-приложений.', 
 '{"Применение": "COBOL широко использовался в банковских и финансовых системах, а также в системах управления данными.",
  "Синтаксис": "COBOL имеет английский синтаксис, что делает его понятным для людей, не знакомых с программированием.",
  "Масштабируемость": "COBOL поддерживает большие объемы данных и может обрабатывать сложные бизнес-логику."}', FALSE);



INSERT INTO projects (user_id, creation_time, status, moderator_id, count) VALUES
(2, NOW(), 'deleted', 1, 2),
(2, NOW(), 'created', 1, 3),
(3, NOW(), 'draft', 1, 0),
(2, NOW(), 'draft', NULL, 0);




INSERT INTO files (lang_id, project_id, code) VALUES
(1, 1, 
'def hello_world():
    print("Hello, World!")

if __name__ == "__main__":
    hello_world()'),

(2, 1, 
'#include <iostream>

using namespace std;

int main() {
    cout << "Hello, World!" << endl;
    return 0;
}'),

(3, 2, 
'package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}'),

(4, 2, 
'<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hello, World!</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>'),

(5, 2, 
'body {
    background-color: lightblue;
}

h1 {
    color: white;
    text-align: center;
}');

