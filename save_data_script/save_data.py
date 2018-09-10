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
def element_insert(filename):
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

# link insert
def link_insert(links_file, elements_file):
    try:
        with open(links_file, encoding='utf-8') as f1:
            links = json.load(f1)
    except IOError as e:
        print(e)
        return
    try:
        with open(elements_file, encoding='utf-8') as f2:
            elements = json.load(f2)
    except IOError as e:
        print(e)
        return

    for link in links:
        scopes = list(link)
        if scopes[0] == scopes[1]:
            continue
        scope_id1 = scopes_map[scopes[0]]
        scope_id2 = scopes_map[scopes[1]]
        scope_name1 = scopes[0]
        scope_name2 = scopes[1]
        element_name1 = link[scopes[0]]
        element_name2 = link[scopes[1]]
        if scope_id1 > scope_id2:
            scope_id1, scope_id2 = scope_id2, scope_id1
            scope_name1, scope_name2 = scope_name2, scope_name1
            element_name1, element_name2 = element_name2, element_name1
        element_code1 = elements[scope_name1][element_name1][0]
        element_code2 = elements[scope_name2][element_name2][0]
        # print(scope_id1, scope_id2, element_code1, element_code2)
        try:
            with connection.cursor() as cursor:
                sql = "INSERT INTO `link` (`scope_id1`, `scope_id2`, `element_code1`, `element_code2`, `status`) VALUES (%s, %s, %s, %s, %s)"
                cursor.execute(sql, (scope_id1, scope_id2, element_code1, element_code2, 0))
                connection.commit()
        except pymysql.err.IntegrityError as e:
            print(e)
            continue


def main():
    # scope_init('scopes.json')
    # element_insert('elements.json')
    link_insert('links.json', 'elements.json')
    connection.close()

if __name__ == "__main__":
    main()


