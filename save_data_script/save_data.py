import pymysql
import json

connection = pymysql.connect(host='localhost', port=3306, user='admin', passwd='mypass', db='questionair')
scopes_map = {'电网工程': 1, '建筑工程': 2, '铁路工程': 3, '公路工程': 4, '水利工程': 5, '民航工程': 6, '石油管道': 7, '城市轨道': 8, '机械制造': 9}

# scope init
def scope_init(filename):
    try:
        with open(filename, encoding='utf-8') as f:
            scopes = json.load(f)
        with connection.cursor() as cursor:
            sql = "INSERT INTO `scope` (`name`, `code`) VALUES (%s, %s)"
            for scope in scopes:
                # print(scope['name'], scope['code'])
                cursor.execute(sql, (scope['name'], scope['code']))
                connection.commit()
    except IOError as e:
        print(e)
        return

# element insert
def element_init(filename):
    try:
        with open(filename, encoding='utf-8') as f:
            elements = json.load(f)
    except IOError as e:
        print(e)
        return
    for scope, datas in elements.items():
        id = scopes_map[scope]
        for name, codes in datas.items():
            # print(id, name, codes[0])
            try:
                with connection.cursor() as cursor:
                    sql = "INSERT INTO `element` (`scope_id`, `name`, `code`) VALUES (%s, %s, %s)"
                    cursor.execute(sql, (id, name, codes[0]))
                    connection.commit()
            except pymysql.err.IntegrityError as e:
                print(e)
                continue

def main():
    # scope_init('scopes.json')
    element_init('elements.json')
    connection.close()

if __name__ == "__main__":
    main()

