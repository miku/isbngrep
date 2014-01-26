#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <pcre.h>
#include <stdbool.h>

#define OVECCOUNT 120

char *str_replace(char *orig, char *rep, char *with) {
    char *result;
    char *ins;
    char *tmp;
    int len_rep, len_with, len_front, count;

    if (!orig) {
        return NULL;
    }
    if (!rep) {
        rep = "";
    }
    if (!with) {
        with = "";
    }
    len_rep = strlen(rep);
    len_with = strlen(with);
    ins = orig;
    for (count = 0; (tmp = strstr(ins, rep)); ++count) {
        ins = tmp + len_rep;
    }
    tmp = result = malloc(strlen(orig) + (len_with - len_rep) * count + 1);
    if (!result) {
        return NULL;
    }
    while (count--) {
        ins = strstr(orig, rep);
        len_front = ins - orig;
        tmp = strncpy(tmp, orig, len_front) + len_front;
        tmp = strcpy(tmp, with) + len_with;
        orig += len_front + len_rep;
    }
    strcpy(tmp, orig);
    return result;
}

bool is_valid_isbn_10(char *s) {
    if (strlen(s) != 10) {
        return false;
    }
    int i = 0, value = 0, sum = 0;
    char bit[2];
    for(i = 0; s[i]; i++) {
        memcpy(bit, &s[i], 1);
        bit[1] = '\0';
        if (strcmp("X", bit) == 0) {
            value = 10;
        } else {
            value = atoi(bit);
        }
        // printf("i = %d, v = %d, sum = %d, bit = %s\n", i, value, sum, bit);
        sum += (10 - i) * value;
    }
    // printf("%s -> %d\n", s, (sum % 11 == 0));
    return (sum % 11 == 0);
}

bool is_valid_isbn_13(char *s) {
    if (strlen(s) != 13) { return false; }
    return true;
}

int main(int argc, char const *argv[])
{
    pcre *re;
    const char *error;
    int erroffset;
    int ovector[OVECCOUNT];
    int rc, i;

    re = pcre_compile(
        "[X0-9-]{10,25}",
        0,
        &error,
        &erroffset,
        NULL);

    if (re == NULL) {
        printf("PCRE compilation failed at offset %d: %s\n",
               erroffset, error);
        return 1;
    }

    char *line = NULL;
    size_t len = 0;
    ssize_t read;

    while ((read = getline(&line, &len, stdin)) != -1) {
        rc = pcre_exec(
            re,
            NULL,
            line,
            (int)strlen(line),
            0,
            0,
            ovector,
            OVECCOUNT);
        if (rc < 0) {
            continue;
        }
        if (rc == 0) {
            rc = OVECCOUNT / 3;
        }
        for (i = 0; i < rc; i++) {
            char *ss_start = line + ovector[2 * i];
            int ss_length = ovector[2 * i + 1] - ovector[2 * i];
            char substring[ss_length];
            sprintf(substring, "%.*s", ss_length, ss_start);
            char *stripped = str_replace(substring, "-", "");
            // printf("candidate: %s\n", stripped);
            if (is_valid_isbn_10(stripped)) {
                printf("%s\n", stripped);    
            }
            free(stripped);
        }
    }
    free(line);
    return 0;
}